/******************************************************************************
* FILENAME:      utils_test.go
*
* AUTHORS:       Qiu Siyu START DATE: 周三 7月 13 2022
*
* LAST MODIFIED: 星期三, 七月 13th 2022, 上午10：40
*
* CONTACT:       siyu.qiu@smartmore.com
******************************************************************************/

package utils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMD5Hash(t *testing.T) {
	cases := []struct {
		Name string
		Pwd  string
		Salt string
	}{
		{"xrw", "123456", "asd919a1sd651s"},
		{"admin", "smartmore2022", "asd919a1sd651s"},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			println(MD5Hash(c.Pwd, c.Salt))
		})
	}
}

func TestRandomCode(t *testing.T) {
	cases := []struct {
		Name string
	}{
		{"xrw"},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			println(RandomCode())
		})
	}
}

func TestRandomString(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			println(RandomString(tt.args.size))
		})
	}
}

func TestOther(t *testing.T) {
	t.Run("timestamp", func(t *testing.T) {
		println(timestamppb.New(time.Now()))
	})
	t.Run("time.Time", func(t *testing.T) {
		fmt.Printf("%v\n", time.Now())
	})
	t.Run("GB", func(t *testing.T) {
		fmt.Printf("%d\n", int64(2<<30))
	})
}

func TestGrpc(t *testing.T) {
	logHelper := log.NewHelper(log.With(log.GetLogger(), "module", "test"))
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial("localhost:52001", opts...)
	if err != nil {
		logHelper.Fatal(err.Error())
	}
	logHelper.Info("success")
	err = conn.Invoke(context.Background(), "/grpc.health.v1.Health/Check", nil, nil)
	if err != nil {
		logHelper.Fatal(err.Error())
		return
	}
}
