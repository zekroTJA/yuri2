/** @format */

import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';
import { HttpClient } from 'selenium-webdriver/http';

/** @format */

@Injectable({
  providedIn: 'root',
})
export class APIService {
  private rootURL = '';

  constructor(private http: HttpClient) {
    this.rootURL = environment.production ? '' : 'http://localhost:8080';
  }
}
