import React from "react";
import PropTypes from "prop-types";
import "./DataPoint.css";

function DataPoint(props) {
  function handleClick(e) {
    e.preventDefault();
    console.log('The link was clicked => ' + props.name);
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
