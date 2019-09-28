/** @format */

import { Component } from '@angular/core';

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.sass'],
})
export class SideBarComponent {
  public sortByName: boolean;
  public inChannel: boolean;

  public isAdmin;
}
