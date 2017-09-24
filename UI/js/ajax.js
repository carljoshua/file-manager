function ajaxObj(method, url){
	var x = new XMLHttpRequest();
	x.open(method, url, true);
	x.setRequestHeader('Content-type', 'application/json');
	return x;
}

function ajaxReturn(x){
	if(x.readyState == 4 && x.status == 200){
		return true;
	}
}

function _(id){
	return document.getElementById(id)
}
