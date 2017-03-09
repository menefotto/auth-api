package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	url := "http://104.155.30.36:8080/api/v1/users/me"
	var wg sync.WaitGroup

	for i := 0; i < 300; i++ {
		wg.Add(1)
		go Req(i, url, &wg)
	}

	wg.Wait()
}

func Req(n int, url string, wg *sync.WaitGroup) {
	defer wg.Done()

	t1 := time.Now()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CRSF-TOKEN", "bwS5y_tvd3plvbSYPhsNayGIAbM:1489075538356")

	req.AddCookie(&http.Cookie{Name: "token", Value: "id=MTQ4OTA3NTUzOHwtbG80TENYTXlfUUdNcFVzb3prYkpVRGVOV29aSWdGeEhNYmtsdUZtT0FocTZ6WDUxSi1TenFuclZQN3NDejJ5Z0ZMUnhYU1hodElOZG1RQ1ZjN3Aybk9jTDBlQkFHMW1TZjlDMDJmRUZxU2dDYy00N2xSeTlPbl94cmg1LUZ2dFlUVHdJU3hGNXNJeVA4enFnS2xDTmtBb3NBdV9OZGtoSUVJOGdYbEt2b0R3b1dLRk9tVDZNVW5tTHRHUDJhR3RIRkxfb3NhOUF1dlRMb1JkdmVqdnFZOG9xZkZ0czNQUnV0dFRTVHpUX0tzQjRqaTVKbVM5WFBma2doRUZYS2ZyU2xFREFWVFhMSGJCanRkREJQLVFteC1DVW9Fa20zRjFLTkhqM2d2QWFROEVDNEJIanY2azFWdlZpTTNPdHd5WThHODN1S190eGV0ak4tdXM3WTVZSFNaaXF4RG5FN2hyY2VxWFJMSUgtOVVSZVh1VjEzRFgxOGRMbkllRnpRUGhCUnpBWTNLOV9jQUlDN2hxN19yVmU5a0UtUW9PSzlkNm5KdHpaaExtazFKNVlOaW4zdVpkMXBOU3pyQzNxSFFxTjRGb2FOLS1EbGpGLXZKRUx4TmJqVnUxUFI1U1lEZ2F5eHpZczZmcUtKbFhTcms2SHZEYUFySGF6VEd3NWZ1eE9wMFc4bUVEa0ZhaTZYaHItalB5b2lqcy1kQkp2NXByMC1ZaGlOVEJHM1BqaGo1SnpRVXR5QVduQjBGS1pROGE4REtvQ0hEcFVOZllHcUpzRHRDcXJweFh2YnFXZ19UWmtjYXBpeFpjSFJxTUd1NFJUZnotNlJBcU1hVlpiLXRtOXRrdmpWZU5hSVBzTzB2bGVUS1h5N0liWXV4OGg1Z3dudElldm1hQmgtZFNIWmlkWUtvWFY5dmdoRGoxYjVEUWVMUT18M6u_DtmYI-1zHd_pDGQRUbiN4ruxJq7NDdLkNyZYWnM="})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	t2 := time.Now()

	fmt.Println("Req number: ", n, "Status: ", resp.Status, "Time taken: ", t2.Sub(t1))
}
