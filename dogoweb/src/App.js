import React, { Component } from 'react';
import './App.css';
import {Switch, Route, BrowserRouter} from 'react-router-dom';
import $ from 'jquery';

//Components
import NavbarGo from './components/NavbarGo';
import HomeScreen from './components/HomeScreen';
import Credit from './components/Credit';


class App extends Component {
  constructor(props){
    super(props);

  }
  componentDidMount(){
    $(document).ready(function(){ 
      if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition((position) => {
          console.log("Latitude: " + position.coords.latitude + 
          "\nLongitude: " + position.coords.longitude); 

          $.ajax({url:'https://maps.googleapis.com/maps/api/geocode/json?latlng='+position.coords.latitude+','+position.coords.longitude+'&key=AIzaSyCMdcgwIyXKwmLtxJQfMT1NFhj1khc5G18', method: 'GET'} ).done(function(response){console.log(response);});

        });
      } else {
        console.log("Geolocation is not supported by this browser.");
      }
    });

  }

  render() {
    return (

      <div>
        <NavbarGo/>
        <HomeScreen/>


        <Credit />
      </div>
      
    );
  }
}

export default App;
