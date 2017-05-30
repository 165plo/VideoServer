package main

import (
    "fmt"
    //"html"
    "log"
    "net/http"
	"io/ioutil"
	"strings"
)

var (
	head string = `<!DOCTYPE html>
<html>
<head>
<title>Page Title</title>
</head>
<body>

<h1>This is a Heading</h1>`

	tail string =	`</body></html>`
	
	script string = `<script>
	var vids = document.getElementsByTagName("video");
	var firstVid = true;
	window.onload = function(){
		for (var i =0;i<vids.length;i++){
			vids[i].load();
			vids[i].oncanplaythrough = function(){
				if(firstVid){
					this.play();
					this.pause();
				}
			}
			
			vids[i].onended = function(){
				var sourceTag = document.getElementsByTagName("source");
				var filename = sourceTag[0].getAttribute("src");
				var newUrl = document.URL.replace(escape(filename)+';',"");
				window.history.pushState("","",newUrl);
				if(!newUrl.endsWith('=')){
					var newVid = newUrl.split('?')[1].split(';')[0].split('=')[1];
					sourceTag[0].src = newVid;
					this.load();
					this.play();
					firstVid = false;
				}
			}
			
		}
	}
	
	function createVidPlaylist(){
		var divs = document.getElementsByTagName("div");
		var vidList = [];
		for (var i=0;i<divs.length;i++){
			var inputInfo = divs[i].getElementsByTagName("input");
			var value = inputInfo[0].value;
			if(value != ""){
				var vid = divs[i].getElementsByTagName("a")[0].getAttribute("href").replace("/Watch?Vids=","");
				var o = {"order":value, "vid":vid};
				vidList.push(o);
			}
		}
		vidList.sort(function sortVids(ob1, ob2){
			return ob1.order - ob2.order;
		});
		
		var videoList = "";
		for (var i=0;i<vidList.length;i++){
			videoList += vidList[i].vid +';';
		}
		document.location.href = document.location.origin+"/Watch?Vids=" + videoList;
	}
	
	
	</script>`
)

func main() {
	http.HandleFunc("/List/", List)
	http.HandleFunc("/Watch", Watch)
	http.Handle("/", http.FileServer(http.Dir("./Vids/")))
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func List(w http.ResponseWriter, r *http.Request){
	folderPath := "./Vids"
	folderPath += r.URL.Path[len("/List"):]
	fmt.Printf("%q\n", r.URL.Path)
	fmt.Printf("%q\n", folderPath)
	fmt.Fprintf(w, head + getVideos(folderPath) + script + tail)
}

func Watch(w http.ResponseWriter, r *http.Request){
	//fileToWatch := r.URL.Path[len("/Watch"):]
	//fmt.Printf("%q", r.URL.Path)
	queryMap := r.URL.Query();
	fileToWatch := queryMap["Vids"][0];
	fileText := ("<video width=\"400\" preload=\"none\" controls><source src=\""+fileToWatch+"\" type=\"video/mp4\">Your browser does not support HTML5 video.</video>")
	
	fmt.Fprintf(w, head + fileText + script + tail)
}

func getVideos(path string) string{
	fileText:=""
	files, err := ioutil.ReadDir(path)
	if(err != nil){
		fileText = `<p>No Files Found</p>`
	}
	for _, file := range files {
		fileText += "<div>"
		fileText += "<input type=\"text\" name=\"order\" size=1>   "
		if strings.Contains(strings.ToLower(file.Name()), "mp4"){
			fileText += ("<a href=\"/Watch?Vids="+ strings.Replace(path, "./Vids", "", -1) + "/" +file.Name()+"\">"+file.Name()+"</a>\n<br/>\n")
		} else if file.IsDir() {
			fileText += ("<a href=\"" + strings.Replace(path, "./Vids", "/List", -1) + "/" + file.Name() + "\">"+file.Name()+"</a>\n<br/>\n")
		} else{
			fileText += ("<label>" + file.Name() + "</label>\n<br/>\n")
		}
		fileText += "</div>"
	}
	fileText += "<button onclick=\"createVidPlaylist()\" type=\"button\">Play Videos</button>"
	return fileText;
}
