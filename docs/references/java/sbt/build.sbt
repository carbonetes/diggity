import Dependencies._

lazy val root = (project in file("."))
  .settings(
    inThisBuild(List(
      organization := "com.example",
      scalaVersion := "2.12.13",
      version := "0.1.0-SNAPSHOT"
    )),
    name := "SparkExample",
    libraryDependencies ++= Seq(
      "org.scala-lang" % "scala-library" % "2.13.6",
      "com.typesafe.akka" %% "akka-actor" % "2.6.17",
      "org.apache.spark" %% "spark-core" % "3.0.1",
      "org.apache.spark" %% "spark-sql" % "3.0.1",
      "org.scalatest" %% "scalatest" % "3.2.9" % "test",
      "org.slf4j" % "slf4j-simple" % "1.7.32"
    ),
    initialCommands in console := """
      import org.apache.log4j.{Level, Logger}
      import org.apache.spark.sql.SparkSession
      import org.apache.spark.sql.functions._
      Logger.getLogger("org.apache.spark").setLevel(Level.WARN)
      val spark = SparkSession.builder
        .master("local[*]")
        .appName("spark-shell")
        .getOrCreate
      import spark.implicits._
      lazy val sc = spark.sparkContext
    """,
    cleanupCommands in console := "if (spark != null) spark.stop()"
  )
