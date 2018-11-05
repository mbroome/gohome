import React from "react";
import PropTypes from "prop-types";

// import the DataPoint component
import DataPoint from "./DataPoint";

function DataPointList(props) {
  return (
    <div>{props.datapoints.map(c => <DataPoint key={c.id} name={c.name} value={c.value} timestamp={c.timestamp} group={c.group} />)}</div>
  );
}

DataPointList.propTypes = {
  datapoints: PropTypes.array.isRequired
};

export default DataPointList;
