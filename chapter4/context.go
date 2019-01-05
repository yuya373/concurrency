package chapter4

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func GreetingAndFarewell() {
	var wg sync.WaitGroup

	done := make(chan interface{})
	// defer close(done)

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printGreeting(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printFarewell(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		close(done)
	}()

	wg.Wait()
}

func printGreeting(done <-chan interface{}) error {
	greeting, err := genGreeting(done)

	if err != nil {
		return err
	}

	fmt.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell(done <-chan interface{}) error {
	farewell, err := genFarewell(done)

	if err != nil {
		return err
	}

	fmt.Printf("%s world!\n", farewell)
	return nil
}

func genGreeting(done <-chan interface{}) (string, error) {
	locale, err := locale(done)

	switch {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(done <-chan interface{}) (string, error) {
	locale, err := locale(done)

	switch {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}

	return "", fmt.Errorf("unsupported locale")
}

func locale(done <-chan interface{}) (string, error) {
	select {
	case <-done:
		return "", fmt.Errorf("canceled")
	case <-time.After(1 * time.Minute):
	}

	return "EN/US", nil
}

func GreetingAndFarewellWithContext() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	locale := func(ctx context.Context) (string, error) {
		if deadline, ok := ctx.Deadline(); ok {
			if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
				// witout calling `time.After(1 * time.Minute)`
				return "", context.DeadlineExceeded
			}
		}

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(1 * time.Minute):
		}
		return "EN/US", nil
	}

	genGreeting := func(ctx context.Context) (string, error) {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		locale, err := locale(ctx)

		switch {
		case err != nil:
			return "", err
		case locale == "EN/US":
			return "hello", nil
		}
		return "", fmt.Errorf("unsupported locale")
	}

	printGreeting := func(ctx context.Context) error {
		greeting, err := genGreeting(ctx)

		if err != nil {
			return err
		}

		fmt.Printf("%s world!\n", greeting)
		return nil
	}

	genFarewell := func(ctx context.Context) (string, error) {
		locale, err := locale(ctx)

		switch {
		case err != nil:
			return "", err
		case locale == "EN/US":
			return "hello", nil
		}

		return "", fmt.Errorf("unsupported locale")
	}

	printFarewell := func(ctx context.Context) error {
		farewell, err := genFarewell(ctx)

		if err != nil {
			return err
		}

		fmt.Printf("%s world!\n", farewell)
		return nil
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printGreeting(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v\n", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printFarewell(ctx); err != nil {
			fmt.Printf("cannot print farewell: %v\n", err)
		}
	}()

	wg.Wait()
}

func DataBag() {
	handleResponse := func(ctx context.Context) {
		fmt.Printf(
			"handling response for %v (%v)",
			ctx.Value("userID"),
			ctx.Value("authToken"),
		)
	}

	processRequest := func(userID, authToken string) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "authToken", authToken)

		handleResponse(ctx)
	}

	processRequest("jane", "abc123")
}

func SafeMapKey() {
	type foo int
	type bar int

	m := make(map[interface{}]int)

	m[foo(1)] = 1
	m[bar(1)] = 2

	fmt.Println(m)
}

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func UserID(c context.Context) string {
	return c.Value(ctxUserID).(string)
}

func AuthToken(c context.Context) string {
	return c.Value(ctxAuthToken).(string)
}

func SafeDataBag() {
	handleResponse := func(c context.Context) {
		fmt.Printf(
			"handling response for %v (auth: %v)",
			UserID(c),
			AuthToken(c),
		)
	}
	processRequest := func(userID, authToken string) {
		ctx := context.WithValue(context.Background(), ctxUserID, userID)
		ctx = context.WithValue(ctx, ctxAuthToken, authToken)

		handleResponse(ctx)
	}

	processRequest("jane", "abc123")
}
