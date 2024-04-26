package frequency_limit_golang

import (
	"context"
	"fmt"
	"testing"
)

func TestIncrAndCheck(t *testing.T) {
	fmt.Println(IncrAndCheck(context.Background(), 1))
	fmt.Println(IncrAndCheck(context.Background(), 1))
	fmt.Println(IncrAndCheck(context.Background(), 1))
	fmt.Println(IncrAndCheck(context.Background(), 1))
	fmt.Println(IncrAndCheck(context.Background(), 1))
	fmt.Println("=====================")
	fmt.Println(IncrAndCheck(context.Background(), 1))
	fmt.Println(IncrAndCheck(context.Background(), 1))
	fmt.Println(IncrAndCheck(context.Background(), 1))
}

func TestCheck(t *testing.T) {
	fmt.Println(Check(context.Background(), 1, 0))
}
