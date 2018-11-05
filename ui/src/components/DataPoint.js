import React from "react";
import PropTypes from "prop-types";
import "./DataPoint.css";

function DataPoint(props) {
  return (
    <div className="datapoint">
      <span>{props.name} => {props.value} => {props.group}</span>
    </div>
  );
}

DataPoint.propTypes = {
  name: PropTypes.string.isRequired,
  value: PropTypes.string.isRequired,
  timestamp: PropTypes.string.isRequired
};

export default DataPoint;
