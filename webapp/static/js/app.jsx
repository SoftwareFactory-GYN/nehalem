let App = React.createClass({
    componentWillMount: function () {
    },
    render: function () {

        if (this.loggedIn) {
            return (<LoggedIn/>);
        } else {
            return (<Home/>);
        }
    }
});

let Home = React.createClass({
    render: function () {
        return (
            <div className="container">
                <div className="col-xs-12 jumbotron text-center">
                    <h1>Nehalem</h1>
                    <a className="btn btn-primary btn-lg btn-login btn-block">Sign In</a>
                </div>
            </div>);
    }
});

let LoggedIn = React.createClass({
    render: function () {
        return (
            <div className="col-lg-12">
                <span className="pull-right"><a onClick={this.logout}>Log out</a></span>
                <h2>Logged in</h2>
            </div>);
    }
});

ReactDOM.render(<App />, document.getElementById('app'));