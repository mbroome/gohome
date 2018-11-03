import React, { Component } from "react";
import "./App.css";

import axios from "axios";

import ContactList from "./components/ContactList";

class App extends Component {
  // default state object
  state = {
    contacts: []
  };

  //async loadData() {
  loadData() {
    axios
      //.get("https://jsonplaceholder.typicode.com/users")
      .get("http://192.168.1.80:8080/topic/")
      .then(response => {
        // create an array of contacts only with relevant data
        const newContacts = response.data.map(c => {
          return {
            id: c.id,
            name: c.name,
            value: c.value,
            timestamp: c.timestamp
          };
        });

        // create a new "state" object without mutating
        // the original state object.
        const newState = Object.assign({}, this.state, {
          contacts: newContacts
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
        <ContactList contacts={this.state.contacts} />
      </div>
    );
  }
}

export default App;
