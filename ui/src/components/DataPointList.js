import React from "react";
import PropTypes from "prop-types";

// import the DataPoint component
import DataPoint from "./DataPoint";

function DataPointList(props) {
  return (
    <div>{props.datapoints.map(c => <DataPoint key={c.id} id={c.id} name={c.name} value={c.value} timestamp={c.timestamp} group={c.group} type={c.type}/>)}</div>
  );
}

DataPointList.propTypes = {
  datapoints: PropTypes.array.isRequired
};

export default DataPointList;
