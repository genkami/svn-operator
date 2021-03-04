package svnconfig_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/genkami/svn-operator/svnconfig"
)

var _ = Describe("Svnconfig", func() {
	Describe("AuthzSVNAccessFile", func() {
		var config *svnconfig.Config
		render := func() string {
			result, err := config.AuthzSVNAccessFile()
			Expect(err).NotTo(HaveOccurred())
			return result
		}

		Describe("section [groups]", func() {
			Context("when the config is empty", func() {
				It("generates an empty config file", func() {
					config = &svnconfig.Config{
						Repositories: []*svnconfig.Repository{},
						Groups:       []*svnconfig.Group{},
						Users:        []*svnconfig.User{},
					}
					Expect(render()).To(Equal(`
[groups]

`))
				})
			})

			Context("when the config contains one group", func() {
				var theGroup string
				BeforeEach(func() {
					theGroup = "engen1"
					config = &svnconfig.Config{
						Repositories: []*svnconfig.Repository{},
						Groups:       []*svnconfig.Group{{Name: theGroup}},
						Users:        []*svnconfig.User{},
					}
				})

				Context("when the groups has no user", func() {
					It("generates an empty group", func() {
						config.Groups[0].Users = []string{}
						Expect(render()).To(Equal(`
[groups]
engen1 = 

`))
					})
				})

				Context("when the group has exactly one user", func() {
					It("generates a group that consists of a single user", func() {
						config.Groups[0].Users = []string{"gura"}
						Expect(render()).To(Equal(`
[groups]
engen1 = gura

`))
					})
				})

				Context("when the group has more than one users", func() {
					It("generates a comma-separated list of users", func() {
						config.Groups[0].Users = []string{"gura", "ame", "ina", "calli", "kiara"}
						Expect(render()).To(Equal(`
[groups]
engen1 = gura, ame, ina, calli, kiara

`))
					})
				})
			})

			Context("when the config contains more than one groups", func() {
				It("generates a list of groups separated by newline", func() {
					config = &svnconfig.Config{
						Repositories: []*svnconfig.Repository{},
						Users:        []*svnconfig.User{},
						Groups: []*svnconfig.Group{
							{"gen4", []string{"coco", "watame", "kanata", "luna", "towa"}},
							{"gen5", []string{"nene", "polka", "lamy", "botan"}},
							{"gen999", []string{}},
						},
					}
					Expect(render()).To(Equal(`
[groups]
gen4 = coco, watame, kanata, luna, towa
gen5 = nene, polka, lamy, botan
gen999 = 

`))
				})
			})
		})
	})
})
