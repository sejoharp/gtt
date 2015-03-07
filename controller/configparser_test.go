package controller

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configparser", func() {
	It("should parse a config.", func() {
		configFile := []byte(`{
  			"tokenKey": "foobarComplex",
  			"salt": "foobar",
			"enableRegister": false,
  			"mongodb": {
    			"host": "localhost",
    			"port": 8888,
    			"database": "test",
    			"user": "user",
    			"password": "password"
  			}
		}`)

		config, err := parseConfig(configFile)

		Expect(err).To(Succeed())
		Expect(config.TokenKey).To(Equal("foobarComplex"))
		Expect(config.Salt).To(Equal("foobar"))
		Expect(config.EnableRegister).To(BeFalse())
		Expect(config.MongoDb.Host).To(Equal("localhost"))
		Expect(config.MongoDb.Port).To(Equal(8888))
		Expect(config.MongoDb.Database).To(Equal("test"))
		Expect(config.MongoDb.User).To(Equal("user"))
		Expect(config.MongoDb.Password).To(Equal("password"))
	})

	It("should read a file.", func() {
		_, err := readFile("../config.json")

		Expect(err).To(Succeed())
	})

	It("should return an error when file reading fails.", func() {
		_, err := readFile("config.json")

		Expect(err).To(HaveOccurred())
	})

	It("should return an error when config reading fails.", func() {
		_, err := ReadConfig("config.json")

		Expect(err).To(HaveOccurred())
	})

	It("should return parsed configfile.", func() {
		config, err := ReadConfig("../config.json")

		Expect(err).To(Succeed())
		Expect(config.TokenKey).To(Equal("foobarComplex"))
		Expect(config.Salt).To(Equal("foobar"))
		Expect(config.EnableRegister).To(BeFalse())
		Expect(config.MongoDb.Host).To(Equal("localhost"))
		Expect(config.MongoDb.Port).To(Equal(8888))
		Expect(config.MongoDb.Database).To(Equal("test"))
		Expect(config.MongoDb.User).To(Equal("user"))
		Expect(config.MongoDb.Password).To(Equal("password"))
	})
})
