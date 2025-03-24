package tests

import (
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("forwarding KumoMTA Log Records", func() {
	When("the Log Record is handled", func() {
		var res *httptest.ResponseRecorder

		JustBeforeEach(func(ctx SpecContext) {
			res = ExecuteRequest(ctx, "./testdata/records/record.json")
			Expect(res.Code).To(Equal(http.StatusAccepted), res.Body.String())
		})

		Specify("the Log Record is forwarded to JetStream", func() {
			Eventually(Must(consumer.FetchNoWait(10)).Messages).
				Should(Receive(HaveField("Data()", MatchJSON(Must(os.ReadFile("./testdata/records/record.json"))))))
		})
	})
})
