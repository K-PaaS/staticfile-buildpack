package integration_test

import (
	"github.com/cloudfoundry/libbuildpack/cutlass"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deploy a staticfile app", func() {
	var app *cutlass.App
	var app_name string
	AfterEach(func() {
		if app != nil {
			app.Destroy()
		}
		app = nil
		app_name = ""
	})
	JustBeforeEach(func() {
		Expect(app_name).ToNot(BeEmpty())
		app = cutlass.New(Fixtures(app_name))
		PushAppAndConfirm(app)
	})

	Context("Using ENV Variable", func() {
		BeforeEach(func() { app_name = "with_https" })

		It("receives a 301 redirect to https", func() {
			_, headers, err := app.Get("/", map[string]string{"NoFollow": "true"})
			Expect(err).To(BeNil())
			Expect(headers).To(HaveKeyWithValue("StatusCode", []string{"301"}))
			Expect(headers).To(HaveKeyWithValue("Location", ConsistOf(HavePrefix("https://"))))
		})

		It("injects x-forwarded-host into Location on redirect", func() {
			var upstreamHostName = "upstreamHostName.com"
			_, headers, err := app.Get("/", map[string]string{"NoFollow": "true", "X-Forwarded-Host": upstreamHostName})
			Expect(err).To(BeNil())
			Expect(headers).To(HaveKeyWithValue("StatusCode", []string{"301"}))
			Expect(headers).To(HaveKeyWithValue("Location", ConsistOf(HavePrefix(fmt.Sprintf("https://%s", upstreamHostName)))))
		})

		Context("Comma separated values in X-Forwarded headers", func() {
			It("picks leftmost x-forwarded-host,-port values into Location on redirect", func() {
				_, headers, err := app.Get("/path1/path2", map[string]string{
					"NoFollow":           "true",
					"X-Forwarded-Host":   "host.com, something.else",
					"X-Forwarded-Prefix": "/pre/fix1, /pre/fix2",
				})
				Expect(err).To(BeNil())
				Expect(headers).To(HaveKeyWithValue("StatusCode", []string{"301"}))
				Expect(headers).To(HaveKeyWithValue("Location", ConsistOf("https://host.com/pre/fix1/path1/path2")))
			})
		})
	})

	Context("Using Staticfile", func() {
		BeforeEach(func() { app_name = "with_https_in_staticfile" })

		It("receives a 301 redirect to https", func() {
			_, headers, err := app.Get("/", map[string]string{"NoFollow": "true"})
			Expect(err).To(BeNil())
			Expect(headers).To(HaveKeyWithValue("StatusCode", []string{"301"}))
			Expect(headers).To(HaveKeyWithValue("Location", ConsistOf(HavePrefix("https://"))))
		})

		It("injects x-forwarded-host into Location on redirect", func() {
			var upstreamHostName = "upstreamHostName.com"
			_, headers, err := app.Get("/", map[string]string{"NoFollow": "true", "X-Forwarded-Host": upstreamHostName})
			Expect(err).To(BeNil())
			Expect(headers).To(HaveKeyWithValue("StatusCode", []string{"301"}))
			Expect(headers).To(HaveKeyWithValue("Location", ConsistOf(HavePrefix(fmt.Sprintf("https://%s", upstreamHostName)))))
		})

		Context("Comma separated values in X-Forwarded headers", func() {
			It("picks leftmost x-forwarded-host,-port values into Location on redirect", func() {
				_, headers, err := app.Get("/path1/path2", map[string]string{
					"NoFollow":           "true",
					"X-Forwarded-Host":   "host.com, something.else",
					"X-Forwarded-Prefix": "/pre/fix1, /pre/fix2",
				})
				Expect(err).To(BeNil())
				Expect(headers).To(HaveKeyWithValue("StatusCode", []string{"301"}))
				Expect(headers).To(HaveKeyWithValue("Location", ConsistOf("https://host.com/pre/fix1/path1/path2")))
			})
		})
	})
})
