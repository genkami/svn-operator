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
				BeforeEach(func() {
					config = &svnconfig.Config{
						Repositories: []*svnconfig.Repository{},
						Groups:       []*svnconfig.Group{{Name: "engen1"}},
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

		Describe("section [REPO_NAME:/]", func() {
			Context("when no group has access to the repository", func() {
				It("drops all permissions", func() {
					config = &svnconfig.Config{
						Repositories: []*svnconfig.Repository{
							{"therepo", []*svnconfig.Permission{}}},
						Groups: []*svnconfig.Group{
							{"fams", []string{"fubuki", "ayame", "mio", "subaru"}}},
						Users: []*svnconfig.User{},
					}
					Expect(render()).To(Equal(`
[groups]
fams = fubuki, ayame, mio, subaru
[therepo:/]
* = 

`))
				})
			})

			Context("when the group have 'r' permission", func() {
				It("grants 'r' permission to the group", func() {
					config = &svnconfig.Config{
						Repositories: []*svnconfig.Repository{
							{"therepo", []*svnconfig.Permission{
								{"smok", "r"},
							}}},
						Groups: []*svnconfig.Group{
							{"smok", []string{"subaru", "mio", "okayu", "korone"}}},
						Users: []*svnconfig.User{},
					}
					Expect(render()).To(Equal(`
[groups]
smok = subaru, mio, okayu, korone
[therepo:/]
* = 
smok = r

`))
				})
			})

			Context("when the group have 'rw' permission", func() {
				It("grants 'rw' permission to the group", func() {
					config = &svnconfig.Config{
						Repositories: []*svnconfig.Repository{
							{"therepo", []*svnconfig.Permission{
								{"idgen2", "rw"},
							}}},
						Groups: []*svnconfig.Group{
							{"idgen2", []string{"ollie", "anya", "reine"}}},
						Users: []*svnconfig.User{},
					}
					Expect(render()).To(Equal(`
[groups]
idgen2 = ollie, anya, reine
[therepo:/]
* = 
idgen2 = rw

`))
				})
			})

			Context("when the permission of the group is explicitly dropped", func() {
				It("grants no permission to the group", func() {
					config = &svnconfig.Config{
						Repositories: []*svnconfig.Repository{
							{"therepo", []*svnconfig.Permission{
								{"nenes", ""},
							}}},
						Groups: []*svnconfig.Group{
							{"nenes", []string{"nenechi", "supernenechi", "hypernenechi"}}},
						Users: []*svnconfig.User{},
					}
					Expect(render()).To(Equal(`
[groups]
nenes = nenechi, supernenechi, hypernenechi
[therepo:/]
* = 
nenes = 

`))
				})
			})

			Context("when more than one groups can have access to the repository", func() {
				It("grants corresponding permissions respectively", func() {
					config = &svnconfig.Config{
						Repositories: []*svnconfig.Repository{
							{"therepo", []*svnconfig.Permission{
								{"board", "r"},
								{"mountains", "rw"},
							}}},
						Groups: []*svnconfig.Group{
							{"board", []string{"shion", "rushia", "kanata", "gura"}},
							{"mountains", []string{"choco", "noel", "coco"}}},
						Users: []*svnconfig.User{},
					}
					Expect(render()).To(Equal(`
[groups]
board = shion, rushia, kanata, gura
mountains = choco, noel, coco
[therepo:/]
* = 
board = r
mountains = rw

`))
				})
			})

			Context("when more than one repositories are defined", func() {
				It("generates list of repositories and its permissions", func() {
					config = &svnconfig.Config{
						Repositories: []*svnconfig.Repository{
							{"therepo1", []*svnconfig.Permission{
								{"edible", "r"},
							}},
							{"therepo2", []*svnconfig.Permission{
								{"edible", "rw"},
								{"carnivore", "r"},
							}},
							{"therepo3", []*svnconfig.Permission{
								{"edible", ""},
								{"carnivore", "r"},
							}},
							{"therepo4", []*svnconfig.Permission{
								{"carnivore", "rw"},
							}},
						},
						Groups: []*svnconfig.Group{
							{"edible", []string{"watame", "ina", "kiara"}},
							{"carnivore", []string{"botan", "gura"}}},
						Users: []*svnconfig.User{},
					}
					Expect(render()).To(Equal(`
[groups]
edible = watame, ina, kiara
carnivore = botan, gura
[therepo1:/]
* = 
edible = r
[therepo2:/]
* = 
edible = rw
carnivore = r
[therepo3:/]
* = 
edible = 
carnivore = r
[therepo4:/]
* = 
carnivore = rw

`))
				})
			})
		})
	})
})
