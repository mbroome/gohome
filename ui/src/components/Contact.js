import React from "react";
import PropTypes from "prop-types";
import "./Contact.css";

function Contact(props) {
  return (
    <div className="contact">
      <span>{props.name} => {props.value} => {props.timestamp} </span>
    </div>
  );
}

Contact.propTypes = {
  name: PropTypes.string.isRequired,
  value: PropTypes.string.isRequired,
  timestamp: PropTypes.string.isRequired
};

export default Contact;
