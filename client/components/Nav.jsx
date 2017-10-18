import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router';

// import backgroundImage from '../now-ui-kit/assets/img/blurred-image-1.jpg';

function Nav({ name, teamName, isFetching, error }) {
  let userBlock;
  if (!isFetching && !error) {
    userBlock = (<span className="nav-link">
      <i className="now-ui-icons users_single-02" /> { name } ({teamName})
    </span>);
  }

  return (
    <nav className="navbar navbar-expand-lg bg-primary fixed-top">
      <div className="container">
        <div className="navbar-translate">
          <a className="navbar-brand" href="/" rel="tooltip" title="Lunch Bot. Lunch, by a bot." data-placement="bottom">
            Lunch Bot
          </a>
          <button className="navbar-toggler navbar-toggler" type="button" data-toggle="collapse" data-target="#navigation" aria-controls="navigation-index" aria-expanded="false" aria-label="Toggle navigation">
            <span className="navbar-toggler-bar bar1" />
            <span className="navbar-toggler-bar bar2" />
            <span className="navbar-toggler-bar bar3" />
          </button>
        </div>

        <div className="collapse navbar-collapse justify-content-end" id="navigation" data-nav-image="/static/assets/img/blurred-image-1.jpg">
          <ul className="navbar-nav mr-auto">
            <li className="nav-item active">
              <Link className="nav-link" to="/manage">
                <i className="now-ui-icons shopping_shop" />
              Places <span className="sr-only">(current)</span>
              </Link>
            </li>
          </ul>
          <ul className="navbar-nav ml-auto justify-content-end">
            <li className="nav-item active">
              {userBlock}
            </li>
            <li className="nav-item">
              <a className="nav-link" href="/manage/api/logout">
                <i className="now-ui-icons media-1_button-power" /> Logout
              </a>
            </li>
          </ul>
        </div>
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
