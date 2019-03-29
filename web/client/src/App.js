import React, { Component } from 'react';
import './App.css';
import { Navbar, NavbarBrand, Row, Col, Modal, ModalHeader, ModalBody, ModalFooter, Button} from 'reactstrap';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.min.css';  
import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faQrcode } from '@fortawesome/free-solid-svg-icons'
import Moment from 'react-moment';
import QRCode from 'qrcode.react';

library.add(faQrcode)
class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      sensors: [],
      activeSensor: "",
      readingsLoading: false,
      readings:[],
      proveReading: null
    };
    
    // Sorry
    this.appStoreUrl = "https://www.youtube.com/watch?v=oHg5SJYRHA0";
    this.playStoreUrl = "https://www.youtube.com/watch?v=oHg5SJYRHA0";

    this.baseUrl = "";
    if(window.location.host === "localhost:3000") {
      this.baseUrl = "http://localhost:8001/";
    }
    this.reloadSensors = this.reloadSensors.bind(this);
    this.reloadReadings = this.reloadReadings.bind(this);
    this.renderSensor = this.renderSensor.bind(this);
    this.showSensor = this.showSensor.bind(this);
    this.toggleProof = this.toggleProof.bind(this);
  }

  componentDidMount() {
    this.reloadSensors()
  }

  reloadSensors() {
    fetch(this.baseUrl + "sensors").then((res) => res.json())
    .then((res) => {
      this.setState({sensors:res});
    }, (err) => {
      toast.error("Something went wrong fetching the sensors")
    });
  }

  reloadReadings() {
    this.setState({readings: [], readingsLoading: true}, () => {
      fetch(this.baseUrl + "readings/" + this.state.activeSensor).then((res) => res.json())
      .then((res) => {
        res.sort((a,b) => b.SensorTimestamp - a.SensorTimestamp);

        this.setState({readings:res, readingsLoading: false});
      }, (err) => {
        toast.error("Something went wrong fetching the readings")
      });
    })
    
  }

  toggleProof() {
    this.setState({proveReading:null});
  }

  renderSensor() {
    var sensor = this.state.sensors.find((v) => v.id === this.state.activeSensor);

    var readings;

    if(this.state.readingsLoading === true) {
      readings = ( <p>Loading...</p>);
    } else {
      readings = this.state.readings.map((r, i) => {
        return (<Row key={i}>
          <Col xs={12} className="text-left">
            <Row>
            <Col xs={10}>
            <b>{r.Statement}</b><br/>
            <Moment fromNow>{new Date(r.SensorTimestamp * 1000)}</Moment>
            </Col>
            <Col xs={2}>
            <Button onClick={((e) => { this.setState({proveReading:r})})}><FontAwesomeIcon icon="qrcode" title="Show proof"/></Button>
            </Col>
            
            </Row>
            <hr/>
          </Col>  
        </Row>);
      });
    }

   

    return ( <Row>
      <Col xs={0} sm={1} md={2} lg={4}>&nbsp;</Col>
      <Col xs={12} sm={10} md={8} lg={4}>
      <Button onClick={((e) => { this.setState({activeSensor:""})})}>&lt; Back</Button>
      <h3>{sensor.name}</h3>
      {readings}
      <Modal isOpen={this.state.proveReading !== null} toggle={this.toggleProof} className={this.props.className}>
          <ModalHeader>Statement proof</ModalHeader>
          <ModalBody>
            <Col xs={12} className="text-center">
            <b>{this.state.proveReading === null ? "" : this.state.proveReading.Statement}</b><br/>
            <Moment fromNow>{new Date(this.state.proveReading  === null ? 0 : this.state.proveReading.SensorTimestamp * 1000)}</Moment>
            <p>&nbsp;</p>
            <QRCode size={256} value={this.state.proveReading  === null ? "blah" : this.state.proveReading.Base64Proof} />
            <p>&nbsp;</p>
            <p>Scan the above QR code using the B_Verify app to independently verify the witnessing of the statement</p>
            <Row>
              <Col xs={6} className="text-center">
                <a rel="noopener noreferrer" target="_blank" href={this.appStoreUrl}><img src="appstore.png" alt="App Store"></img></a>
              </Col>
              <Col xs={6} className="text-center">
                <a rel="noopener noreferrer" target="_blank" href={this.playStoreUrl}><img src="playstore.png" alt="Play Store"></img></a>
              </Col>
            </Row>
            </Col>
          </ModalBody>
          <ModalFooter>
            <Button color="secondary" onClick={this.toggleProof}>Close</Button>
          </ModalFooter>
        </Modal>
      </Col>
      <Col xs={0} sm={1}  md={2} lg={4}>&nbsp;</Col>
    </Row>
    )
  }

  showSensor(id) {
    this.setState({activeSensor:id}, () => {
      this.reloadReadings();
    });
  }

  

  render() {
    var sensors = this.state.sensors.map((s) => {
      return (<Row key={s.id}>
        <Col xs={12}>

          <button className="link-button"  onClick={((e) => { this.showSensor(s.id); })}><b>{s.name}</b></button>
          <p>{s.description}</p>
        </Col>
      </Row>)
    })

    if( this.state.sensors.length == 0 ){
      sensors = (<Row>
        <Col xs={12}>

          <p>Sorry, no sensors available</p>
        </Col>
      </Row>)
    }
    
    var mainComponent;

    if(this.state.activeSensor === '') {
      mainComponent = (<Col>
      <Row>
        <Col xs={0} sm={1} md={2} lg={4}>&nbsp;</Col>
        <Col xs={12} sm={10} md={8} lg={4}>
        <h3>What is this?</h3>
        <p>This is a demonstration of what you can do using <a href="https://bverify.org/">B_Verify</a>. This demonstration shows you secure sensor readings that were witnessed to the Bitcoin blockchain. You can see the list of sensors below this introduction. Each of these sensors has their readings witnessed by B_Verify, and posts the proofs for that witnessing to this website. Using the B_Verify Mobile App for <a href={this.appStoreUrl}>iOS</a> and <a href={this.playStoreUrl}>Android</a> you can independently verify these witness statements. Click on one of the sensors to view their data.</p>
      
        </Col>
        <Col xs={0} sm={1}  md={2} lg={4}>&nbsp;</Col>
      </Row>
      <Row>
        <Col xs={0} sm={1} md={2} lg={4}>&nbsp;</Col>
        <Col xs={12} sm={10} md={8} lg={4}>
        <h3>Available sensors:</h3>
        
        {sensors}
      
        </Col>
        <Col xs={0} sm={1}  md={2} lg={4}>&nbsp;</Col>
      </Row>
      
        </Col>);
    } else {
      mainComponent = this.renderSensor()
    }

    return (
      <div className="App">
        <Navbar expand="xs">
            <NavbarBrand className="mx-auto" href="/">
              <span className="nav-link" style={{color:'black'}}><b>Sensors powered by </b><img alt="BVerify" src="logo.png" /></span>
            </NavbarBrand>
        </Navbar>
    {mainComponent}
        <ToastContainer
          position="top-right"
          autoClose={5000}
          hideProgressBar={false}
          newestOnTop={false}
          closeOnClick
          rtl={false}
          pauseOnVisibilityChange={false}
          draggable={false}
          pauseOnHover={false}
        />
      </div>
    );
  }
}

export default App;
