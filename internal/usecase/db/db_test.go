package db

import (
	"context"
	"github.com/22Fariz22/shorturl/internal/entity"
	"github.com/22Fariz22/shorturl/internal/handler"
	repoMock "github.com/22Fariz22/shorturl/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func _TestInDBRepositoryGetURL(t *testing.T) {
	ctl := gomock.NewController(t)
	ctl.Finish()

	repo := repoMock.NewMockRepository(ctl)

	ctx := context.Background()
	short := handler.GenUlid()

	mockResp := entity.URL{LongURL: "https://ya.ru"}

	expected := entity.URL{LongURL: "https://ya.ru"}

	repo.EXPECT().GetURL(ctx, short).Return(mockResp, true).Times(1)
	long, ok := repo.GetURL(ctx, short)

	require.Equal(t, expected, long)
	require.EqualValues(t, true, ok)
}
