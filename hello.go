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
				if(firstVid){
					var sourceTag = document.getElementsByTagName("source");
					sourceTag[0].src = "/mov_bbb.mp4";
					this.load();
					this.play();
					firstVid = false;
				}
			}
			
		}
	}
	</script>`
)

func main() {
	http.HandleFunc("/List", List)
	http.HandleFunc("/Watch/", Watch)
	http.Handle("/", http.FileServer(http.Dir("./Vids/")))
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func List(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, head + getVideos() + tail)
}

func Watch(w http.ResponseWriter, r *http.Request){
	fileToWatch := r.URL.Path[len("/view/"):]
	fileText := ("<video width=\"400\" preload=\"none\" controls><source src=\""+fileToWatch+"\" type=\"video/mp4\">Your browser does not support HTML5 video.</video>")
	fmt.Fprintf(w, head + fileText + script + tail)
}

func getVideos() string{
	fileText:=""
	files, err := ioutil.ReadDir("./Vids")
	if(err != nil){
		fileText = `<p>No Files Found</p>`
	}
	for _, file := range files {
		if strings.Contains(strings.ToLower(file.Name()), "mp4"){
			fileText += ("<a href=\"Watch/"+file.Name()+"\">"+file.Name()+"</a>\n<br/>\n")
		}else{
			fileText += ("<p>" + file.Name() + "</p>\n<br/>\n")
		}
		
	}
	return fileText;
}
