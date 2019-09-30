/** @format */

import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { SideBarComponent } from './components/sidebar/sidebar.component';
import { SliderComponent } from './components/slider/slider.component';
import { FormsModule } from '@angular/forms';
import { ToastComponent } from './components/toast/toast.component';
import { MainRouteComponent } from './routes/main.route';
import { SoundBtnComponent } from './components/soundbtn/soundbtn.component';
import { HttpClientModule } from '@angular/common/http';
import { SearchBarComponent } from './components/searchbar/searchbar.component';
import { ContextMenuComponent } from './components/contextmenu/contextmenu.component';

@NgModule({
  declarations: [
    MainRouteComponent,

    AppComponent,
    SideBarComponent,
    SliderComponent,
    ToastComponent,
    SoundBtnComponent,
    SearchBarComponent,
    ContextMenuComponent,
  ],
  imports: [BrowserModule, AppRoutingModule, FormsModule, HttpClientModule],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
