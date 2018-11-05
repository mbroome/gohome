import React, { Component } from "react";
import { Column, Row } from 'simple-flexbox';
import "./App.css";

import axios from "axios";

import DataPointList from "./components/DataPointList";

class App extends Component {
  // default state object
  state = {
    datapoints: []
  };

  //async loadData() {
  loadData() {
    axios
      //.get("https://jsonplaceholder.typicode.com/users")
      .get("http://192.168.1.80:8080/data")
      .then(response => {
        var now = new Date();
        var newDataPoints = {};
        for(var i = 0, l = response.data.length; i < l; i++) {
           !(response.data[i].group in newDataPoints) && (newDataPoints[response.data[i].group] = [])
           newDataPoints[response.data[i].group].push({
                                                     id: response.data[i].id,
                                                     name: response.data[i].name,
                                                     value: response.data[i].value,
                                                     group: response.data[i].group,
                                                     type: response.data[i].type,
                                                     timestamp: response.data[i].timestamp
                                                   });
        }

        console.log(newDataPoints);

        // create a new "state" object without mutating
        // the original state object.
        const newState = Object.assign({}, this.state, {
           datapoints: newDataPoints,
           polltime: now.toLocaleString()
        });

        // store the new state object in the component's state
        this.setState(newState);
      })
      .catch(error => console.log(error));
  };

  componentDidMount() {
    this.loadData()
    setInterval(() => this.loadData(), 30000);
  }

  render() {
    return (
      <div className="App">
      <Column flexGrow={1}>
          <Row horizontal='center'>
              <h1>{this.state.polltime}</h1>
          </Row>
          <Row vertical='start'>
             {Object.entries(this.state.datapoints).map(([key, value]) => (
                <Column flexGrow={1} key={key} horizontal='center'>
                  <div className="datapoint">{key}</div>
                  <DataPointList key={key} datapoints={value} />
                </Column>
             ))}
             <Column flexGrow={1} horizontal='center'>
                <div className="datapoint"><span>logs</span></div>
             </Column>
          </Row>
      </Column>
      </div>

    );
  }
}

export default App;
