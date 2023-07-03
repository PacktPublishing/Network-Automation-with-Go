# Network Automation with Go

<a href="https://www.packtpub.com/product/network-automation-with-go/9781800560925?utm_source=github&utm_medium=repository&utm_campaign=9781800560925"><img src="https://static.packt-cdn.com/products/9781800560925/cover/smaller" alt="Early Access" height="256px" align="right"></a>

This is the code repository for [Network Automation with Go](https://www.packtpub.com/product/network-automation-with-go/9781800560925?utm_source=github&utm_medium=repository&utm_campaign=9781800560925), published by Packt.

**Learn how to automate network operations and build applications using the Go programming language**

## What is this book about?
Go’s built-in first-class concurrency mechanisms make it an ideal choice for long-lived low-bandwidth I/O operations, which are typical requirements of network automation and network operations applications.

This book covers the following exciting features:
* Understand Go programming language basics via network-related examples
* Find out what features make Go a powerful alternative for network automation
* Explore network automation goals, benefits, and common use cases
* Discover how to interact with network devices using a variety of technologies
* Integrate Go programs into an automation framework
* Take advantage of the OpenConfig ecosystem with Go
* Build distributed and scalable systems for network observability

If you feel this book is for you, get your [copy](https://www.amazon.com/dp/1800560923) today!

<a href="https://www.packtpub.com/?utm_source=github&utm_medium=banner&utm_campaign=GitHubBanner"><img src="https://raw.githubusercontent.com/PacktPublishing/GitHub/master/GitHub.png" 
alt="https://www.packtpub.com/" border="5" /></a>

## Instructions and Navigations
All of the code is organized into folders. For example, ch01.

The code will look like the following:
```
func main() {     
    a := -1     
    
    var b uint32     
    b = 4294967295     
    
    var c float32 = 42.1
}
```

**Following is what you need for this book:**
This book is for all network engineers, administrators, and other network practitioners looking to understand what network automation is and how the Go programming language can help develop network automation solutions. As the first part of the book offers a comprehensive overview of Go’s main features, this book is suitable for beginners with a solid grasp on programming basics.

With the following software and hardware list you can run all code files present in the book.
### Software and Hardware List
| Software required | OS required |
| ------------------------------------ | ----------------------------------- |
| Go 1.18.1 | Linux (Ubuntu 22.04, Fedora 35), Windows Subsystem for Linux (WSL2) or macOS |
| Containerlab 0.28.1 | Linux (Ubuntu 22.04, Fedora 35), Windows Subsystem for Linux (WSL2) or macOS |
| Docker 20.10.14 | Linux (Ubuntu 22.04, Fedora 35), Windows Subsystem for Linux (WSL2) or macOS |


We also provide a PDF file that has color images of the screenshots/diagrams used in this book. [Click here to download it](https://packt.link/hOgov).

## Errata
* Page 64 - In the second code block; the second comment shows the Binary value of the localhost variable which has a decimal value of 127,0,0,1 as 1111 1111,0000 0000,0000 0000,0000 0001 which is not correct and it should be 0111 1111,0000 0000,0000 0000,0000 0001.
* Page 64 - In the second code block; the ipv4 in the first and second comment is meant to be the ipAddr and localhost variables respectively.

### Related products
*  Python Network Programming Techniques[[Packt]](https://www.packtpub.com/product/python-network-programming-techniques/9781838646639?utm_source=github&utm_medium=repository&utm_campaign=9781838646639) [[Amazon]](https://www.amazon.com/dp/1838646639)

*  Network Automation Cookbook[[Packt]](https://www.packtpub.com/product/network-automation-cookbook/9781789956481?utm_source=github&utm_medium=repository&utm_campaign=9781789956481) [[Amazon]](https://www.amazon.com/dp/178995648X)

## Get to Know the Authors
**Nicolas**
is a Specialist Solutions Architect at Red Hat. In his role, he helps customers of all sizes to automate the provisioning and operation of IT infrastructure, services, and applications. Prior to that, he worked in the networking industry for 15 years, becoming a Cisco Certified Design Expert (CCDE) and Internetwork Expert (CCIE). He is passionate about writing open-source software in Go with a keen interest in cloud technologies.


**Michael**
is a Cloud Infrastructure Solutions Architect, currently working in the networking business unit of NVIDIA. Throughout his career, he held multiple roles ranging from network operations, through software development to systems architecture and design. He enjoys breaking the boundaries between different disciplines and coming up with creative solutions to satisfy business needs and solve technical problems in the most optimal way. He is a prolific open-source contributor and writer with much of his work focused on cloud-native infrastructure, automation, and orchestration.

### Download a free PDF

 <i>If you have already purchased a print or Kindle version of this book, you can get a DRM-free PDF version at no cost.<br>Simply click on the link to claim your free PDF.</i>
<p align="center"> <a href="https://packt.link/free-ebook/9781800560925">https://packt.link/free-ebook/9781800560925 </a> </p>
