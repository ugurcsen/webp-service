package main

import (
	"net/http"
	"testing"
	"time"
)

var preapered = false

func prepareForBenchmarkAndTest() {
	if preapered {
		return
	}
	preapered = true
	go main()
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))
	go http.ListenAndServe(":8081", mux)
	time.Sleep(10 * time.Second)
}

func testConvert(t *testing.T, width string) {
	prepareForBenchmarkAndTest()
	resp, err := http.Get("http://localhost:" + port + "/?url=http://localhost:8081/static/img.png&w=" + width)
	if err != nil {
		t.Fatal(err)
	}
	byt := make([]byte, 2<<20)
	resp.Body.Read(byt)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.StatusCode)
	}
}

func TestConverting1080(t *testing.T) {
	testConvert(t, "1080")
}

func TestConverting1080WithQuality(t *testing.T) {
	testConvert(t, "1080")
}

func TestConverting720(t *testing.T) {
	testConvert(t, "720")
}

func TestConverting480(t *testing.T) {
	testConvert(t, "480")
}

func benchmarkConverting(b *testing.B, width string) {
	b.StopTimer()
	prepareForBenchmarkAndTest()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get("http://localhost:" + port + "/?url=http://localhost:8081/static/img.png&w=" + width)
		if err != nil {
			b.Fatal(err)
		}
		byt := make([]byte, 2<<20)
		resp.Body.Read(byt)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			b.Fatal(resp.StatusCode)
		}
	}
}

func BenchmarkConverting1080(b *testing.B) {
	benchmarkConverting(b, "1080")
}

func BenchmarkConverting720(b *testing.B) {
	benchmarkConverting(b, "720")
}

func BenchmarkConverting480(b *testing.B) {
	benchmarkConverting(b, "480")
}

func BenchmarkConverting320(b *testing.B) {
	benchmarkConverting(b, "320")
}

func BenchmarkConverting240(b *testing.B) {
	benchmarkConverting(b, "240")
}

func BenchmarkConverting64(b *testing.B) {
	benchmarkConverting(b, "64")
}

func BenchmarkConverting32(b *testing.B) {
	benchmarkConverting(b, "32")
}
