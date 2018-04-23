package home

func (this *HomeController) HomeIndexAction() {
	this.Tplname="views/admin/index.html"
	this.Rander()
}
