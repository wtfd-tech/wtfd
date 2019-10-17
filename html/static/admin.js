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
    uve.value = res.Name;
    uvd.value = res.DisplayName;
    uvp.value = res.Points;
    uva.checked = res.Admin;
    showDialog(uv);

  });


}
