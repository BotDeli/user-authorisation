package tests

import (
	"AccountControl/internal/storage/redis/session"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

const lifetime = 1 * time.Second

type RedisMock struct {
	mock.Mock
}

func (m *RedisMock) Close() error {
	return nil
}

func (m *RedisMock) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *RedisMock) Get(key string) *redis.StringCmd {
	args := m.Called(key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *RedisMock) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(key, expiration)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *RedisMock) Del(keys ...string) *redis.IntCmd {
	args := m.Called(keys)
	return args.Get(0).(*redis.IntCmd)
}

func TestErrorNewSession(t *testing.T) {
	initMocks := func(client *RedisMock) {
		cmd := redis.NewStatusResult("", testError)
		client.On("Set", mock.Anything, testID, lifetime).Return(cmd)
	}

	testMocks := func(t *testing.T, r session.Redis) {
		_, err := r.NewSession(testID)
		if err == nil {
			t.Error("expected error, got nil")
		}
	}

	testRedisMock(t, initMocks, testMocks)
}

func testRedisMock(t *testing.T, initMocks func(*RedisMock), testMocks func(*testing.T, session.Redis)) {
	client := new(RedisMock)
	initMocks(client)
	r := session.Redis{Client: client, Lifetime: lifetime}
	testMocks(t, r)
}

func TestSuccessfulNewSession(t *testing.T) {
	initMocks := func(client *RedisMock) {
		cmd := redis.NewStatusResult("", nil)
		client.On("Set", mock.Anything, testID, lifetime).Return(cmd)
	}

	testMocks := func(t *testing.T, r session.Redis) {
		_, err := r.NewSession(testID)
		if err != nil {
			t.Error(err)
		}
	}

	testRedisMock(t, initMocks, testMocks)
}

func TestErrorGetEmailFromSession(t *testing.T) {
	initMocks := func(client *RedisMock) {
		cmd := redis.NewStringResult("", testError)
		client.On("Get", testKey).Return(cmd)
	}

	testMocks := func(t *testing.T, r session.Redis) {
		_, err := r.GetIdFromSession(testKey)
		if err == nil {
			t.Error("expected error, got nil")
		}
	}

	testRedisMock(t, initMocks, testMocks)
}

func TestSuccessfulGetEmailFromSession(t *testing.T) {
	initMocks := func(client *RedisMock) {
		cmd := redis.NewStringResult(testID, nil)
		client.On("Get", testKey).Return(cmd)
	}

	testMocks := func(t *testing.T, r session.Redis) {
		id, err := r.GetIdFromSession(testKey)
		if err != nil {
			t.Error(err)
		}
		if id != testID {
			t.Errorf("expected %s, got %s", testID, id)
		}
	}

	testRedisMock(t, initMocks, testMocks)
}

func TestFunctionalityUpdateSessionLifeTime(t *testing.T) {
	initMocks := func(client *RedisMock) {
		cmd := redis.NewBoolResult(true, nil)
		client.On("Expire", testKey, lifetime).Return(cmd).Once()
	}

	testMocks := func(t *testing.T, r session.Redis) {
		r.UpdateSessionLifeTime(testKey)
	}

	testRedisMock(t, initMocks, testMocks)
}
