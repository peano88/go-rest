
function parseResponse(response) {
  document.getElementById("main").innerHTML = response;
}

function main() {
  var xmlhttp = new XMLHttpRequest();
  var url = "/api/notes";

  xmlhttp.onreadystatechange = function() {
      if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
	  var response = JSON.parse(xmlhttp.responseText);
	  parseResponse(response);
      }
  };
  xmlhttp.open("GET", url, true);
  xmlhttp.send();
}
