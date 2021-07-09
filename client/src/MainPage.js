import React, { Component } from "react";
import axios from "axios";
import { Card, Header, Form, Icon, Image, Grid ,Divider} from "semantic-ui-react";
import ReactHtmlParser, { processNodes, convertNodeToElement, htmlparser2 } from 'react-html-parser';

let endpoint = "http://localhost:8888";

class MainPage extends Component {
  constructor(props) {
    super(props)
    this.state = {
      task: "",
      items: []
    };
  }

  componentDidMount() {
    this.getPOST();
  }

  onChange = event => {
    this.setState({
      [event.target.name]: event.target.value
    });
  };

  onSubmit = () => {
    let { task } = this.state;
    // console.log("pRINTING task", this.state.task);
    if (task) {
      axios
        .post(
          endpoint + "/api/v1/post", {

            
          task
        }, {
          headers: {
            "Content-Type": "application/x-www-form-urlencoded"
          }
        }
        )
        .then(res => {
          this.getPOST();
          this.setState({
            task: ""
          });
          console.log(res);
        });
    }
  };

  getPOST = () => {
    axios.get(endpoint + "/api/v1/postall").then(res => {
      console.log(res);
      if (res.data) {
        this.setState({
          items: res.data.map((item, keyss) => {
            let color = "yellow";
            console.log(keyss)
            let picture = "http://localhost:8888/assets/image.png";

            if (item.picture) {
              picture = endpoint + "/upload/" + item.picture;
            }

            if (item.status) {
              color = "green";
            }
            if (keyss % 2 == 0) {
              var rows = 1
            } else {
              var rows = 2
            }

            return (
          

                  <Grid.Column key={item.id} divided='vertically'>
                    <Card  color={color} fluid >
                      <Image src={picture} wrapped ui={false} />
                      <Card.Content>
                        <Card.Header> {item.title} </Card.Header>
                        <Card.Meta>
                          <span className='date'>{item.created_at}   </span>
                        </Card.Meta>
                        <Card.Description>
                          {ReactHtmlParser(item.content)}
                        </Card.Description>
                      </Card.Content>

                      <Card.Content extra>
                        <Grid>
                          <Grid.Column floated='left' width={5}>
                            <a onClick={() => this.updatePost(item.id)} >
                              <Icon onClick={() => this.updatePost(item.id)} name='edit' />Edit </a>
                          </Grid.Column>
                          <Grid.Column floated='left' width={5}>
                            <a onClick={() => this.updatePost(item.id)} >
                              <Icon onClick={() => this.deletePost(item.id)} name='delete' /> Delete </a>
                          </Grid.Column>
                        </Grid>
                      </Card.Content>
                    </Card>
                  </Grid.Column>
              
         

            );


          })
        });
      } else {
        this.setState({
          items: []
        });
      }
    });
  };

  updatePost = id => {
    axios
      .put(endpoint + "/api/v1/task/" + id, {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        }
      })
      .then(res => {
        console.log(res);
        this.getPOST();
      });
  };

  /*
    statusTask = id => {
      axios
        .put(endpoint + "/api/v1/task/" + id, {
          headers: {
            "Content-Type": "application/x-www-form-urlencoded"
          }
        })
        .then(res => {
          console.log(res);
          this.getPOST();
        });
    };
  */
  deletePost = id => {
    axios
      .delete(endpoint + "/api/v1/post/" + id, {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        }
      })
      .then(res => {
        console.log(res);
        this.getPOST();
      });
  };
  render() {
    // 
    return (<div>
      <div className="row" >
        < Header className="header" as="h2" >
          Post Add
          </Header>
      </div>

      <Form onSubmit={this.onSubmit}>
        <Form.Group widths='equal'>
          <Form.Input fluid label='Title' placeholder='Title' />


        </Form.Group>

        <Form.TextArea label='About' placeholder='Tell us more about you...' />
        <Form.Button>Submit</Form.Button>
      </Form>

      <Grid columns={2}>
        <Grid.Row>
           {this.state.items} 
        </Grid.Row>
       </Grid>
       </div>
    );
  }
}

export default MainPage;