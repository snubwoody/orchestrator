mod error;
use reqwest::Client;
use serde::{Serialize,Deserialize};

// TODO check the rules about the name
#[derive(Debug,Serialize,Deserialize,Default)]
#[serde(rename_all="kebab-case")]
pub struct Instance{
    name: String,
    description: String,
    #[serde(rename="ip-forwarding")]
    can_ip_forward: bool,
    machine_type: String,
    disks: Vec<Disk>
}

#[derive(Debug,Serialize,Deserialize)]
#[serde(rename_all="kebab-case")]
pub struct Disk{
    name: String,
    #[serde(rename="type")]
    disk_type: DiskType,
    mode: DiskMode,
    /// Indicates whether this is the boot disk
    boot: bool,
    /// Specifies whether the disk will be delete when the instance is deleted
    auto_delete: bool,
    /// The disk size in GB
    disk_size: i64
}

#[derive(Debug,Serialize,Deserialize)]
struct MachineType{
    zone: String,
    _type: String
}

#[derive(Debug,Serialize,Deserialize)]
#[serde(rename_all="lowercase")]
enum DiskType{
    Scratch,
    Persistent
}

#[derive(Debug,Serialize,Deserialize)]
#[serde(rename_all="kebab-case")]
enum DiskMode{
    ReadWrite,
    ReadOnly
}


async fn insert_instance(project: &str,zone: &str){
    let url = format!(
        "https://compute.googleapis.com/compute/v1/projects/{project}/zones/{zone}/instances"
    );

    let client = Client::new();
    let instance = Instance::default();
    let response = client.post(url).json(&instance).send().await.unwrap();
    let body: serde_json::Value = response.json().await.unwrap();
    dbg!(body);
}


#[cfg(test)]
mod tests{ 
    use super::*;
    use toml::{toml, Value};

    #[tokio::test]
    async fn send_reqwest(){
        insert_instance("", "").await;
    }

    #[test]
    fn parse_instance_toml(){
        let data = toml! {
            name = "log-service"
            description = "Micro-service for logging requests and actions"
            ip-forwarding = true
            machine-type = "zones/us-east-5a/machines/debian-12"

            [[disks]]
            boot = true
            auto-delete = true
            disk-size = 10
            type = "persistent"
            mode = "read-write"
        };

        let instance = toml::from_str::<Instance>(&data.to_string());
    }
}