package e2e_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("data-gather", func() {
	var (
		manifestPath string
		tmpDir       string
	)

	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "ginkgo-run")
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		err := os.RemoveAll(tmpDir)
		Ω(err).ShouldNot(HaveOccurred())
	})

	act := func(manifestPath string) (session *gexec.Session, err error) {
		args := []string{"util", "data-gather", "-m", manifestPath, "-b", "../testing/assets", "-g", "log-api"}
		cmd := exec.Command(cliPath, args...)
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		return
	}

	Context("when manifest exists", func() {
		BeforeEach(func() {
			manifestPath = "../testing/assets/gatherManifest.yml"
		})

		It("gathers data to stdoout", func() {
			session, err := act(manifestPath)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))
			output := session.Out.Contents()
			Expect(output).Should(ContainSubstring(`"properties.yaml":"name: cf`))

			//fmt.Printf("%s\n", string(output))
			var yml map[string]interface{}
			err = json.Unmarshal(output, &yml)
			Expect(err).ToNot(HaveOccurred())

			// for a, _ := range yml {
			//         fmt.Printf("%#v\n", a)
			// }
			fmt.Printf("%s\n", yml["properties.yaml"])

			// var data map[string]interface{}
			// err = json.Unmarshal(yml["properties.yaml"].([]byte), &data)
			// Expect(err).ToNot(HaveOccurred())
			// fmt.Printf("%s\n", data)

		})
	})
})
