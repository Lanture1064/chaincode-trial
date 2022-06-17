package main_test

import (
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {

	var tempDir string
	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
	})

	DescribeTable("Running the build command produces the correct error code",
		func(expectedErrorCode int, getArgs func() []string) {
			args := getArgs()
			command := exec.Command(buildCmdPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(expectedErrorCode))
		},
		Entry("When the image.json file exists", 0, func() []string { return []string{"./testdata/validimage", "CHAINCODE_METADATA_DIR", tempDir} }),
		Entry("When the image.json file does not exist", 1, func() []string { return []string{"CHAINCODE_SOURCE_DIR", "CHAINCODE_METADATA_DIR", "BUILD_OUTPUT_DIR"} }),
		Entry("When the image.json file is invalid", 1, func() []string { return []string{"./testdata/invalidimage", "CHAINCODE_METADATA_DIR", tempDir} }),
		Entry("When too few arguments are provided", 1, func() []string { return []string{"CHAINCODE_SOURCE_DIR"} }),
		Entry("When too many arguments are provided", 1, func() []string {
			return []string{"CHAINCODE_SOURCE_DIR", "CHAINCODE_METADATA_DIR", "BUILD_OUTPUT_DIR", "UNEXPECTED_ARGUMENT"}
		}),
	)

	It("should copy chaincode metadata to the build output directory", func() {
		args := []string{"./testdata/withmetadata", "CHAINCODE_METADATA_DIR", tempDir}
		command := exec.Command(buildCmdPath, args...)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		indexPath := filepath.Join(tempDir, "META-INF", "test", "test.txt")
		Expect(indexPath).To(BeARegularFile())
	})
})
