import dialogPolyfill from "dialog-polyfill";
import {showDialog} from './util';

export default class AdminPage {
  uv: HTMLDialogElement;
  uve: HTMLInputElement;
  uvd: HTMLInputElement;
  uvp: HTMLInputElement;
  uva: HTMLInputElement;
  table: Element;

  constructor() {
    this.uv = <HTMLDialogElement> document.getElementById("ud");
    this.uve = <HTMLInputElement> document.getElementById("namechange");
    this.uvd = <HTMLInputElement> document.getElementById("displaynamechange");
    this.uvp = <HTMLInputElement> document.getElementById("pointschange");
    this.uva = <HTMLInputElement> document.getElementById("adminchange");
    this.table = document.getElementsByClassName("tbody")[0];
    dialogPolyfill.registerDialog(this.uv);
    for (let e of this.table.getElementsByClassName("show-dialog-button")) {
      e.addEventListener("click", () => {
        this.update(e.id);
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
