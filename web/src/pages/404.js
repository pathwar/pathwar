import React from "react";
import { Link } from "gatsby";

const browser = typeof window !== "undefined" && window;

const NotFoundPage = () => {
  return (
    browser && (
      <div>
        <h4>404 Page Not Foundto="/">Go back to homepage</Link>
      </div>
    )
  );
};

export default NotFoundPage;
