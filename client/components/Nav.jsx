import React, { PropTypes } from 'react';
import { Link } from 'react-router';

function Nav() {
  return (
    <nav className="navbar navbar-toggleable-md navbar-inverse fixed-top bg-inverse">
      <button
        className="navbar-toggler navbar-toggler-right hidden-lg-up"
        type="button"
        data-toggle="collapse"
        data-target="#navbarsExampleDefault"
        aria-controls="navbarsExampleDefault"
        aria-expanded="false"
        aria-label="Toggle navigation"
      >
        <span className="navbar-toggler-icon" />
      </button>
      <a className="navbar-brand" href="/">Lunch Bot</a>

      <div className="collapse navbar-collapse" id="navbarsExampleDefault">
        <ul className="navbar-nav mr-auto">
          <li className="nav-item active">
            <Link className="nav-link" to="/manage">
              Places <span className="sr-only">(current)</span>
            </Link>
          </li>
          <li className="nav-item">
            <a className="nav-link" href="/manage/api/logout">Logout</a>
          </li>
        </ul>
      </div>
    </nav>
  );
}

export default Nav;
