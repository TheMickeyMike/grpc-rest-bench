package benchmarks

type JobBuilder struct {
	client  *HTTPClient
	retries int
	target  interface{}
}

func NewJobBuilder() *JobBuilder {
	return &JobBuilder{
		client: NewHTTPClient(HTTP),
	}
}

func (j *JobBuilder) WithHttpClient(client *HTTPClient) *JobBuilder {
	j.client = client
	return j
}

func (j *JobBuilder) WithRetries(count int) *JobBuilder {
	j.retries = count
	return j
}

// func (j *JobBuilder) Build() wpool.ExecutionFn {
// 	return func(ctx context.Context) (string, error) {
// 		var (
// 			err      error
// 			response ResponseDetails
// 			target   pb.UserAccount
// 			retries  int = 3
// 		)
// 		for retries > 0 {
// 			response, err = j.client.MakeRequest(ctx, "https://localhost:8080/api/v1/users/61df07d341ed08ad981c143c", &target)
// 			if err != nil {
// 				retries -= 1
// 			} else {
// 				break
// 			}
// 		}
// 		return response.StatKey(), nil
// 	}
// }
