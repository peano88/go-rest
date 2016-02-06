function get() {
  var xmlhttp = new XMLHttpRequest();
  var url = "/api/notes";

  xmlhttp.onreadystatechange = function() {
      if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
	  var response = JSON.parse(xmlhttp.responseText);
          document.getElementById("d01").innerHTML += response + '<br>';
      }
  };
  xmlhttp.open("GET", url, true);
  xmlhttp.send();
}

function main() {
  var source = new EventSource('/api/notes/events');
  source.onmessage = function(e) {
    document.getElementById("d01").innerHTML += e.data + '<br>';
  };
}
