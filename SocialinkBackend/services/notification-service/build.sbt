name := "socialink-notification-service"

version := "1.0.0"

scalaVersion := "2.13.12"

lazy val akkaVersion = "2.8.5"
lazy val akkaHttpVersion = "10.5.3"
lazy val slickVersion = "3.5.0"

libraryDependencies ++= Seq(
  // Akka
  "com.typesafe.akka" %% "akka-actor-typed" % akkaVersion,
  "com.typesafe.akka" %% "akka-stream" % akkaVersion,
  "com.typesafe.akka" %% "akka-cluster-typed" % akkaVersion,
  "com.typesafe.akka" %% "akka-cluster-sharding-typed" % akkaVersion,
  
  // Akka HTTP
  "com.typesafe.akka" %% "akka-http" % akkaHttpVersion,
  "com.typesafe.akka" %% "akka-http-spray-json" % akkaHttpVersion,
  
  // Database
  "com.typesafe.slick" %% "slick" % slickVersion,
  "com.typesafe.slick" %% "slick-hikaricp" % slickVersion,
  "org.postgresql" % "postgresql" % "42.7.1",
  
  // Redis
  "com.github.etaty" %% "rediscala" % "1.9.0",
  
  // JSON
  "io.spray" %% "spray-json" % "1.3.6",
  "com.typesafe.play" %% "play-json" % "2.10.3",
  
  // Firebase Cloud Messaging (Push notifications)
  "com.google.firebase" % "firebase-admin" % "9.2.0",
  
  // Email
  "com.sun.mail" % "javax.mail" % "1.6.2",
  
  // WebSocket
  "com.typesafe.akka" %% "akka-stream" % akkaVersion,
  
  // Config
  "com.typesafe" % "config" % "1.4.3",
  
  // Logging
  "com.typesafe.akka" %% "akka-slf4j" % akkaVersion,
  "ch.qos.logback" % "logback-classic" % "1.4.14",
  
  // Kafka (for event consumption)
  "com.typesafe.akka" %% "akka-stream-kafka" % "4.0.2",
  
  // Testing
  "com.typesafe.akka" %% "akka-actor-testkit-typed" % akkaVersion % Test,
  "org.scalatest" %% "scalatest" % "3.2.17" % Test
)

// Assembly settings for fat JAR
assembly / assemblyMergeStrategy := {
  case PathList("META-INF", xs @ _*) => MergeStrategy.discard
  case "reference.conf" => MergeStrategy.concat
  case x => MergeStrategy.first
}

assembly / mainClass := Some("com.socialink.notification.Main")
