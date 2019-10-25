var uv = document.getElementById("ud");
var uve = document.getElementById("namechange");
var uvd = document.getElementById("displaynamechange");
var uvp = document.getElementById("pointschange");
var uva = document.getElementById("adminchange");
(function(){
if(location.href.match("admin").length === 1){
    dialogPolyfill.registerDialog(uv);
} 
})();

function update(email){
  fetch('/getUserData/'+email).then( (res) => res.json()).then((res) => {
    console.log(res);
    uve.value = res.name;
    uvd.value = res.displayname;
    uvp.value = res.points;
    uva.checked = res.admin;
    showDialog(uv);

  });


}
