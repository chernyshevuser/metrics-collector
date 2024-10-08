package impl_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/test-go/testify/assert"

	businessimpl "github.com/chernyshevuser/practicum-metrics-collector/internal/server/business/impl"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/handler/impl"
	storageimpl "github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage/impl"
	mocklogger "github.com/chernyshevuser/practicum-metrics-collector/tools/logger/mock"
)

func TestUpdateMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	defer logger.Sync()

	logger.EXPECT().Info("goodbye from business-svc").Times(1)
	logger.EXPECT().Info("goodbye from db-svc").Times(1)
	logger.EXPECT().Sync().Times(1)

	dbSvc, err := storageimpl.New(context.TODO(), logger)
	if err != nil {
		t.Errorf("cant create db svc: %v", err)
	}
	defer dbSvc.Close()

	businessSvc := businessimpl.New(dbSvc, logger)
	defer businessSvc.Close()

	svc := impl.New(businessSvc, logger)

	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		endp string
		want want
	}{
		{
			name: "correct counter",
			endp: "/update/counter/sample_text/123",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "incorrect counter by value",
			endp: "/update/counter/sample_text/123.5",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "correct gauge with int value",
			endp: "/update/gauge/sample_text/123",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "correct gauge with float value",
			endp: "/update/gauge/name/123.4",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "incorrect gauge by value",
			endp: "/update/gauge/sample_text/value",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.endp, nil)
			w := httptest.NewRecorder()
			svc.UpdateMetric(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
