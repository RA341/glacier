package listutils

import "sync"

func ToMap[T any, Q any](input []T, mapper func(T) Q) []Q {
	var result = make([]Q, 0, len(input))
	for _, t := range input {
		result = append(result, mapper(t))
	}
	return result
}

func ToMapErr[T any, Q any](input []T, mapper func(T) (Q, error)) ([]Q, error) {
	var result = make([]Q, 0, len(input))
	for _, t := range input {
		q, err := mapper(t)
		if err != nil {
			return nil, err
		}
		result = append(result, q)
	}
	return result, nil
}

func ParallelLoop[T any, R any](input []R, mapper func(R) (T, bool)) []T {
	contChan := make(chan T, len(input))

	var wg sync.WaitGroup
	for _, cont := range input {
		wg.Add(1)
		go func(i R) {
			defer wg.Done()
			res, ok := mapper(i)
			if ok {
				contChan <- res
			}
		}(cont)
	}

	go func() {
		wg.Wait()
		close(contChan)
	}()

	var result []T
	for c := range contChan {
		result = append(result, c)
	}

	return result
}
