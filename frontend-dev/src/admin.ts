import dialogPolyfill from "dialog-polyfill";

export default class AdminPage {
  uv: Element;
  uve: Element;
  uvd: Element;
  uvp: Element;
  uva: Element;
  table: Element;

  constructor() {
    this.uv = document.getElementById("ud");
    this.uve = document.getElementById("namechange");
    this.uvd = document.getElementById("displaynamechange");
    this.uvp = document.getElementById("pointschange");
    this.uva = document.getElementById("adminchange");
    this.table = document.getElementsByClassName("tbody")[0];
    dialogPolyfill.registerDialog(this.uv);
    for (let e of this.table.getElementsByClassName("show-dialog-button")) {
      e.addEventListener("click", () => {
        update(e.id);
      });
    }
  }

  update(email: string) {
    fetch("/getUserData/" + email)
      .then(res => res.json())
      .then(res => {
        console.log(res);
        this.uve.value = res.name;
        this.uvd.value = res.displayname;
        this.uvp.value = res.points;
        this.uva.checked = res.admin;
        showDialog(this.uv);
      });
  }
}
