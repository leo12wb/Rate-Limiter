package rate_limit

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/leo12wb/Rate-Limiter/internal/entity/web_session"
	"github.com/leo12wb/Rate-Limiter/internal/value_objects"
	"github.com/stretchr/testify/suite"
)

type RateLimitTestSuite struct {
	suite.Suite
	ctx context.Context
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(RateLimitTestSuite))
}

func (suite *RateLimitTestSuite) SetupSuite() {
	suite.ctx = context.Background()
}
func (suite *RateLimitTestSuite) TestRateLimitIP() {
	db, mock := redismock.NewClientMock()

	session, err := web_session.NewWebSession("10.0.0.1", "", value_objects.NewRequestLimit(10, 15, 60))
	suite.NoError(err)
	repo := *NewRateLimitRepositoryRedis(db)
	mock.ExpectGet(session.GetRequestTimerId()).SetVal(fmt.Sprintf("%d", time.Now().Unix()))
	mock.ExpectWatch(session.GetRequestCounterId()).RedisNil()
	mock.ExpectTxPipeline()
	mock.ExpectGet(session.GetRequestCounterId()).SetVal("10")
	mock.ExpectSet(session.GetRequestCounterId(), 9, 0).SetVal("9")
	mock.ExpectTxPipelineExec()

	denied, err := repo.DecreaseTokenBucket(&session)
	suite.NoError(err)
	suite.False(denied)

}
func (suite *RateLimitTestSuite) TestRateLimitIPFailWatch() {
	db, mock := redismock.NewClientMock()

	session, err := web_session.NewWebSession("10.0.0.1", "", value_objects.NewRequestLimit(10, 15, 60))
	suite.NoError(err)
	repo := *NewRateLimitRepositoryRedis(db)
	mock.ExpectGet(session.GetRequestTimerId()).SetVal(fmt.Sprintf("%d", time.Now().Unix()))

	for i := 0; i < 5; i++ {
		mock.ExpectWatch(session.GetRequestCounterId()).SetErr(errors.New("Key Changed"))
	}
	mock.ExpectTxPipeline()
	mock.ExpectGet(session.GetRequestCounterId()).SetVal("10")
	mock.ExpectSet(session.GetRequestCounterId(), 9, 0).SetVal("9")
	mock.ExpectTxPipelineExec()

	denied, err := repo.DecreaseTokenBucket(&session)
	suite.Error(err)
	suite.True(denied)

}
func (suite *RateLimitTestSuite) TestRateLimitIPThrottle() {
	db, mock := redismock.NewClientMock()

	session, err := web_session.NewWebSession("10.0.0.1", "", value_objects.NewRequestLimit(10, 15, 60))
	suite.NoError(err)
	repo := *NewRateLimitRepositoryRedis(db)
	mock.ExpectGet(session.GetRequestTimerId()).SetVal(fmt.Sprintf("%d", time.Now().Unix()))

	//mock.MatchExpectationsInOrder(true)
	mock.ExpectWatch(session.GetRequestCounterId())

	mock.ExpectGet(session.GetRequestCounterId()).SetVal("0")
	//mock.ExpectSet(session.GetRequestCounterId(), -1, 0).SetVal("-1")

	denied, err := repo.DecreaseTokenBucket(&session)
	throttledError := ThrottledError{}
	suite.Equal(throttledError.ThrottledError().Error(), err.Error())
	suite.True(denied)

}
func (suite *RateLimitTestSuite) TestRateLimitAPI() {
	db, mock := redismock.NewClientMock()

	session, err := web_session.NewWebSession("10.0.0.1", "LUCAO", value_objects.NewRequestLimit(10, 15, 60))
	suite.NoError(err)
	repo := *NewRateLimitRepositoryRedis(db)
	mock.ExpectGet(session.GetRequestTimerId()).SetVal(fmt.Sprintf("%d", time.Now().Unix()))

	mock.ExpectWatch(session.GetRequestCounterId()).RedisNil()
	mock.ExpectTxPipeline()
	mock.ExpectGet(session.GetRequestCounterId()).SetVal("10")
	mock.ExpectSet(session.GetRequestCounterId(), 9, 0).SetVal("9")
	mock.ExpectTxPipelineExec()

	denied, err := repo.DecreaseTokenBucket(&session)
	suite.NoError(err)
	suite.False(denied)

}
func (suite *RateLimitTestSuite) TestRateLimitAPIFailWatch() {
	db, mock := redismock.NewClientMock()

	session, err := web_session.NewWebSession("10.0.0.1", "LUCAO", value_objects.NewRequestLimit(10, 15, 60))
	suite.NoError(err)
	repo := *NewRateLimitRepositoryRedis(db)
	mock.ExpectGet(session.GetRequestTimerId()).SetVal(fmt.Sprintf("%d", time.Now().Unix()))

	for i := 0; i < 5; i++ {
		mock.ExpectWatch(session.GetRequestCounterId()).SetErr(errors.New("Key Changed"))
	}
	mock.ExpectTxPipeline()
	mock.ExpectGet(session.GetRequestCounterId()).SetVal("10")
	mock.ExpectSet(session.GetRequestCounterId(), 9, 0).SetVal("9")
	mock.ExpectTxPipelineExec()

	denied, err := repo.DecreaseTokenBucket(&session)
	suite.Error(err)
	suite.True(denied)

}
func (suite *RateLimitTestSuite) TestRateLimitAPIThrottle() {
	db, mock := redismock.NewClientMock()

	session, err := web_session.NewWebSession("10.0.0.1", "LUCAO", value_objects.NewRequestLimit(10, 15, 60))
	suite.NoError(err)
	repo := *NewRateLimitRepositoryRedis(db)
	mock.ExpectGet(session.GetRequestTimerId()).SetVal(fmt.Sprintf("%d", time.Now().Unix()))

	//mock.MatchExpectationsInOrder(true)
	mock.ExpectWatch(session.GetRequestCounterId())

	mock.ExpectGet(session.GetRequestCounterId()).SetVal("0")
	//mock.ExpectSet(session.GetRequestCounterId(), -1, 0).SetVal("-1")

	denied, err := repo.DecreaseTokenBucket(&session)
	throttledError := ThrottledError{}
	suite.Equal(throttledError.ThrottledError().Error(), err.Error())
	suite.True(denied)

}
