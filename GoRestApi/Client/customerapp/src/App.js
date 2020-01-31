import React, { Component } from 'react'
import Customers from './Components/customersc'
import './App.css';

class App extends Component {


  state = {
    customers: [],
    iscust: Boolean,
    eleme: ''
  }

  componentDidMount() {

    fetch('http://localhost:8000/ByScore/')
      .then(res => res.json())
      .then((data) => {
        this.setState({ customers: data })
      }).then((data) => {
        console.log(data)
      })
      .catch(console.log)
  }
  




  render() {
    return (
      <div>
        <h6> besht score</h6>
        <Customers clients={this.state.customers} />
      </div>
    )
  }
}

export default App;
