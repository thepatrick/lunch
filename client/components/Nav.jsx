import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router';

function Nav({ name, teamName, isFetching, error }) {
  let userBlock;
  if (!isFetching && !error) {
    userBlock = (<span className="nav-link">
      { name } ({teamName})
    </span>);
  }

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
        </ul>
        <ul className="navbar-nav ml-auto justify-content-end">
          <li className="nav-item active">
            {userBlock}
          </li>
          <li className="nav-item">
            <a className="nav-link" href="/manage/api/logout">Logout</a>
          </li>
        </ul>
      </div>
    </nav>
  );
}

Nav.propTypes = {
  name: PropTypes.string.isRequired,
  teamName: PropTypes.string.isRequired,
  isFetching: PropTypes.bool.isRequired,
  error: PropTypes.instanceOf(Error),
};

Nav.defaultProps = {
  error: undefined,
};

export default Nav;
