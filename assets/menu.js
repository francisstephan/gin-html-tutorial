function setmenu(option) {
    // get the origin of the url
    // https://stackoverflow.com/questions/40558641/change-url-and-reload-page-with-that-new-url-using-javascript-jquery
    var host = window.location.origin;
    window.location.replace(host + option );
}

function listalbum() {
  var retVal = prompt("Album Id : ", "type Id here");
  if(retVal!="") setmenu("/albums/"+retVal);
}

function updatealbum() {
  var retVal = prompt("Id of album to update :", "type Id here");
  if(retVal!="") setmenu("/update/"+retVal);
}
