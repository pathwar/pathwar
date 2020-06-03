import React from 'react';
import { Link } from "gatsby";

const browser = typeof window !== "undefined" && window;

const NotFoundPage = () => {
  return (
    browser && (
      <div>
        <h4>
          404 Page Not Found
        </h4>
        <Link to="/">Go back to homepage</Link>
      </div>
    )
  );
};

export default NotFoundPage;
