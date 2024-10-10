package fx

import (
	"github.com/gowool/theme"

	"github.com/gowool/cms"
	"github.com/gowool/cms/repository"
	cmstheme "github.com/gowool/cms/theme"
)

func FuncMap(pageRepo repository.Page, menu cms.Menu, matcher cms.Matcher) theme.FuncMap {
	return cmstheme.NewFuncMap(pageRepo, menu, matcher).FuncMap
}
