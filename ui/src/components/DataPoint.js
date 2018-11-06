import React from "react";
import PropTypes from "prop-types";
import axios from "axios";
import "./DataPoint.css";

function DataPoint(props) {
  function handleClick(e) {
    e.preventDefault();
    console.log(props);

    if (props.type === 'switch'){
       var nv = props.value === 'ON'? 'OFF' : 'ON';
       console.log('The link was clicked => ' + props.id + ' => ' + nv);

       axios
         .put("http://192.168.1.80:8080/command/" + props.id, nv)
         .then(response => {
           console.log(response);

         })
         .catch(error => console.log(error));


    }
  }

  return (
    <div className="datapoint">
      <span onClick={handleClick}>{props.name} => {props.value}</span>
    </div>
  );
}

DataPoint.propTypes = {
  name: PropTypes.string.isRequired,
  value: PropTypes.string.isRequired,
  timestamp: PropTypes.string.isRequired
};

export default DataPoint;
